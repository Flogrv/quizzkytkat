# 🚀 Déploiement sur Coolify

Guide pour déployer le quiz de cybersécurité sur Coolify.

## Prérequis

- Coolify installé sur votre VPS
- Domaine configuré (ex: `quizz.yantekc.com`)
- Repository Git accessible

## Étapes de déploiement

### 1. Configuration DNS

Ajoutez un enregistrement DNS pour votre sous-domaine :

```
Type: A
Nom: quizz
Valeur: [IP de votre VPS]
TTL: 3600
```

Ou directement :
```
quizz.yantekc.com -> [IP VPS]
```

### 2. Création du projet dans Coolify

1. Connectez-vous à Coolify
2. Créez un nouveau projet
3. Sélectionnez "New Resource" > "Public Repository" ou "Private Repository"
4. Entrez l'URL de votre repo Git

### 3. Configuration du service

Dans les paramètres du service :

#### General Settings
- **Name**: `cybersec-quiz`
- **Build Pack**: `Dockerfile`

#### Ports & Networking
- **Port**: `2222`
- **Protocol**: `TCP`
- **Public**: `Yes`

Si vous voulez utiliser votre domaine avec un reverse proxy SSH, vous devrez configurer :
- **Domain**: `quizz.yantekc.com`
- **Port externe**: `2222`

**Note**: Pour SSH, il est recommandé d'exposer directement le port sans reverse proxy HTTP.

#### Volumes
Ajoutez ces volumes persistants :

```
./data:/app/data
./questions.json:/app/questions.json
```

Ou via l'interface Coolify :
1. Cliquez sur "Add Volume"
2. Source: `/data` (local sur le VPS)
3. Destination: `/app/data`
4. Répétez pour `questions.json`

#### Environment Variables
Variables optionnelles (si vous modifiez le code pour les utiliser) :
```
SSH_PORT=2222
DB_PATH=/app/data/quiz.db
QUESTIONS_PATH=/app/questions.json
TZ=Europe/Paris
```

### 4. Build Settings

Coolify détectera automatiquement le `Dockerfile`. Vérifiez :
- **Dockerfile Location**: `./Dockerfile`
- **Build Context**: `/`

### 5. Déploiement

1. Cliquez sur "Deploy"
2. Attendez la fin du build
3. Vérifiez les logs

### 6. Configuration du firewall

Sur votre VPS, assurez-vous que le port 2222 est ouvert :

```bash
# UFW
sudo ufw allow 2222/tcp

# iptables
sudo iptables -A INPUT -p tcp --dport 2222 -j ACCEPT
```

### 7. Test de connexion

Testez depuis votre machine locale :

```bash
ssh -p 2222 quizz.yantekc.com
```

Ou avec l'IP directement :
```bash
ssh -p 2222 [IP_VPS]
```

## Configuration avancée

### Utiliser le port 22 standard

Si vous voulez utiliser le port 22 pour le quiz :

1. Changez le port SSH système sur le VPS (ex: 2222)
2. Configurez le quiz pour écouter sur le port 22
3. Mappez le port dans Coolify : `22:2222`

**⚠️ Attention**: Ne verrouillez pas votre accès SSH au VPS !

### Multiple instances avec load balancing

Pour gérer beaucoup de connexions :

1. Créez plusieurs instances du service
2. Utilisez un volume partagé pour la DB
3. Configurez un load balancer TCP devant

### Backup automatique

Ajoutez un script de backup pour la DB :

```bash
#!/bin/bash
# backup.sh
DATE=$(date +%Y%m%d_%H%M%S)
cp /data/quiz.db /backups/quiz_${DATE}.db
# Garder seulement les 7 derniers backups
ls -t /backups/quiz_*.db | tail -n +8 | xargs rm -f
```

Configurez une cron job dans Coolify ou sur le VPS.

### Monitoring

Pour surveiller votre service :

1. Logs en temps réel :
```bash
docker logs -f [container_name]
```

2. Dans Coolify : Section "Logs"

3. Métriques : CPU, RAM, connexions actives

## Mise à jour

Pour mettre à jour l'application :

1. Push votre nouveau code sur Git
2. Dans Coolify : "Redeploy"
3. Coolify rebuild automatiquement l'image

Ou avec webhook auto-deploy :
- Configurez le webhook dans Coolify
- Ajoutez-le dans votre repo Git (GitHub/GitLab)
- Chaque push déploiera automatiquement

## Troubleshooting

### Port déjà utilisé
```bash
# Vérifier quel process utilise le port 2222
sudo netstat -tulpn | grep 2222
sudo lsof -i :2222
```

### Impossible de se connecter
1. Vérifiez que le conteneur tourne : `docker ps`
2. Vérifiez les logs : `docker logs [container_name]`
3. Vérifiez le firewall : `sudo ufw status`
4. Testez depuis le VPS : `ssh -p 2222 localhost`

### Base de données verrouillée
Si SQLite est verrouillé :
```bash
# Arrêtez le service
docker stop [container_name]
# Supprimez le lock
rm /data/quiz.db-shm /data/quiz.db-wal
# Redémarrez
docker start [container_name]
```

### Questions non chargées
Vérifiez que `questions.json` est bien monté :
```bash
docker exec [container_name] cat /app/questions.json
```

## Sécurité

### Limiter les connexions

Ajoutez fail2ban pour limiter les tentatives :

```bash
# /etc/fail2ban/jail.local
[ssh-quiz]
enabled = true
port = 2222
filter = sshd
logpath = /var/log/quiz.log
maxretry = 5
bantime = 3600
```

### Authentification par clé

Pour forcer l'authentification par clé SSH, modifiez `main.go` pour ajouter une validation des clés publiques.

### Rate limiting

Ajoutez un rate limiting au niveau du serveur SSH dans le code pour limiter les connexions par IP.

## Support

En cas de problème :
1. Consultez les logs Coolify
2. Vérifiez les logs du conteneur
3. Testez en local d'abord

## Ressources

- [Documentation Coolify](https://coolify.io/docs)
- [SSH Protocol](https://www.ssh.com/academy/ssh/protocol)
- [Go Wish Documentation](https://github.com/charmbracelet/wish)

---

**Bon déploiement ! 🎉**
