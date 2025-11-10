# 🔐 Cybersec Quiz - SSH Interactive Quiz

Application de quiz de cybersécurité accessible via SSH, construite avec Go et Bubbletea.

## ✨ Fonctionnalités

- 🎯 **Quiz interactif** avec questions par catégories
- 🏆 **Leaderboard** global et par catégorie
- 👤 **Authentification** par pseudo
- 🎨 **Interface TUI** moderne et colorée avec Bubbletea
- 📊 **Suivi des scores** avec persistance SQLite
- 🐳 **Containerisé** pour déploiement facile
- 🔒 **Accès SSH** sécurisé

## 🚀 Démarrage rapide

### Prérequis

- Go 1.22+
- Docker (optionnel)
- SQLite

### Installation locale

1. Clone le repository :
```bash
git clone <your-repo>
cd quizz_cybersec
```

2. Installe les dépendances :
```bash
go mod download
```

3. Lance l'application :
```bash
go run main.go
```

4. Connecte-toi depuis un autre terminal :
```bash
ssh -p 2222 localhost
```

## 🐳 Déploiement Docker

### Build et run avec Docker Compose

```bash
docker-compose up -d
```

### Build manuel

```bash
docker build -t cybersec-quiz .
docker run -p 2222:2222 -v ./data:/app/data -v ./questions.json:/app/questions.json cybersec-quiz
```

## 🎮 Utilisation

1. Connecte-toi via SSH :
```bash
ssh -p 2222 quizz.yantekc.com
```

2. Entre ton pseudo (3 caractères minimum)

3. Choisis une option dans le menu :
   - 🎯 Jouer - Toutes les questions
   - 📚 Jouer - Par catégorie
   - 🏆 Leaderboard Global
   - 📊 Leaderboard Par Catégorie
   - 🚪 Quitter

4. Réponds aux questions avec les flèches ↑/↓ (ou j/k) et valide avec Enter

## 📝 Configuration des questions

Les questions sont stockées dans `questions.json` avec le format suivant :

```json
{
  "questions": [
    {
      "id": 1,
      "category": "Réseau",
      "text": "Quel protocole est utilisé pour sécuriser HTTP ?",
      "options": [
        "SSL/TLS",
        "FTP",
        "SMTP",
        "DNS"
      ],
      "answer": 0
    }
  ]
}
```

- `id` : Identifiant unique
- `category` : Catégorie de la question
- `text` : Texte de la question
- `options` : Liste des réponses possibles
- `answer` : Index de la bonne réponse (commence à 0)

## 🗂️ Structure du projet

```
.
├── main.go              # Point d'entrée, serveur SSH
├── models/              # Structures de données
│   └── models.go
├── storage/             # Gestion DB et questions
│   ├── database.go
│   └── questions.go
├── ui/                  # Interfaces Bubbletea
│   ├── styles.go        # Styles et couleurs
│   ├── username.go      # Écran de connexion
│   ├── menu.go          # Menu principal
│   ├── category_select.go
│   ├── quiz.go          # Interface de quiz
│   └── leaderboard.go   # Affichage des scores
├── questions.json       # Questions du quiz
├── data/                # Base de données SQLite (créé auto)
├── Dockerfile           # Image Docker
└── docker-compose.yml   # Configuration Docker Compose
```

## 🎨 UI/UX

L'interface utilise une palette de couleurs moderne :
- 🟢 Primaire : `#00ff9f` (vert cyan)
- 🟣 Secondaire : `#7d56f4` (violet)
- 🟡 Accent : `#ff6ac1` (rose)
- 🔴 Erreur : `#ff4757` (rouge)
- ✅ Succès : `#2ed573` (vert)

Navigation :
- `↑`/`↓` ou `j`/`k` : Naviguer
- `Enter` : Valider
- `q` ou `Ctrl+C` : Quitter

## 🏗️ Déploiement sur VPS avec Coolify

1. Configure ton DNS :
   - Crée un enregistrement A pour `quizz.yantekc.com` → IP de ton VPS

2. Dans Coolify :
   - Crée un nouveau projet
   - Utilise le Dockerfile
   - Configure le port mapping : `2222:2222`
   - Monte les volumes :
     - `./data:/app/data`
     - `./questions.json:/app/questions.json`

3. Déploie et connecte-toi :
```bash
ssh -p 2222 quizz.yantekc.com
```

## 📊 Base de données

SQLite est utilisé pour stocker :
- Les scores des utilisateurs
- Les statistiques (parties jouées, joueurs uniques)
- L'historique des tentatives

La DB est créée automatiquement au premier lancement dans `./data/quiz.db`.

## 🔒 Sécurité

- Les clés SSH host sont générées automatiquement au premier lancement
- Pas d'authentification stricte par défaut (accessible à tous)
- Les pseudos sont libres (pas de compte)
- Pour ajouter l'authentification par clé SSH publique, modifie le middleware wish

## 🛠️ Développement

### Ajouter des catégories

Il suffit d'ajouter des questions avec de nouvelles valeurs de `category` dans `questions.json`.

### Modifier les styles

Les styles sont centralisés dans `ui/styles.go`.

### Ajouter des fonctionnalités

1. Crée un nouveau modèle dans `ui/`
2. Ajoute un état dans `appModel` (main.go)
3. Implémente `Init()`, `Update()`, et `View()`

## 📦 Dépendances principales

- `github.com/charmbracelet/bubbletea` - Framework TUI
- `github.com/charmbracelet/lipgloss` - Styles terminal
- `github.com/charmbracelet/wish` - Serveur SSH
- `github.com/charmbracelet/bubbles` - Composants UI
- `github.com/mattn/go-sqlite3` - Driver SQLite

## 📄 Licence

MIT

## 🤝 Contribution

Les contributions sont bienvenues ! N'hésite pas à ouvrir une issue ou une PR.

## 👨‍💻 Auteur

Créé pour un projet Epitech - Quiz de cybersécurité

---

**Enjoy & Happy Hacking! 🚀🔐**
# quizzkytkat
