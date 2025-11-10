#!/bin/bash

# Script de déploiement pour VPS/Coolify
# Usage: ./deploy.sh

set -e

echo "🚀 Démarrage du déploiement..."

# Vérifier que Docker est installé
if ! command -v docker &> /dev/null; then
    echo "❌ Docker n'est pas installé. Installez Docker et réessayez."
    exit 1
fi

# Vérifier que docker-compose est disponible
if ! command -v docker-compose &> /dev/null; then
    if ! docker compose version &> /dev/null; then
        echo "❌ docker-compose n'est pas installé. Installez docker-compose et réessayez."
        exit 1
    fi
    DOCKER_COMPOSE="docker compose"
else
    DOCKER_COMPOSE="docker-compose"
fi

# Créer le dossier data si nécessaire
echo "📁 Création du dossier data..."
mkdir -p data

# Vérifier que questions.json existe
if [ ! -f "questions.json" ]; then
    echo "⚠️  questions.json n'existe pas, copie de l'exemple..."
    cp questions.example.json questions.json
fi

# Build l'image Docker
echo "🐳 Build de l'image Docker..."
$DOCKER_COMPOSE build

# Arrêter les anciens conteneurs
echo "🛑 Arrêt des anciens conteneurs..."
$DOCKER_COMPOSE down || true

# Démarrer le nouveau conteneur
echo "🚀 Démarrage du conteneur..."
$DOCKER_COMPOSE up -d

# Attendre que le serveur soit prêt
echo "⏳ Attente du démarrage du serveur..."
sleep 3

# Vérifier les logs
echo "📋 Logs du serveur:"
$DOCKER_COMPOSE logs --tail=20

echo ""
echo "✅ Déploiement terminé !"
echo "🔗 Connectez-vous avec: ssh -p 2222 localhost"
echo "📊 Voir les logs: $DOCKER_COMPOSE logs -f"
echo "🛑 Arrêter le serveur: $DOCKER_COMPOSE down"
