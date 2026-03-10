# chitest

Ce projet est une application web écrite en Go utilisant le routeur [chi](https://github.com/go-chi/chi) et le moteur de templates HTML natif. Elle propose :

- Un serveur HTTP avec gestion des routes et des templates
- Un middleware qui log les headers HTTP en entrée et en sortie
- Un middleware qui ajoute un identifiant unique à chaque requête (header `X-Request-ID`)
- Un fragment HTML dynamique pour la page hello, compatible avec htmx
- Un support du hot reload avec [Air](https://github.com/air-verse/air)

## Prérequis
- Go 1.23+
- [Air](https://github.com/air-verse/air) pour le développement (hot reload)

## Installation

```sh
git clone https://github.com/<votre-utilisateur>/chitest.git
cd chitest
go mod tidy
```

## Lancer en développement

```sh
air
```

## Lancer les tests

```sh
go test ./...
```

## Structure du projet

- `main.go` : point d'entrée, configuration du serveur et des middlewares
- `main_test.go` : tests unitaires
- `static/` : fichiers statiques (CSS, JS, etc.)
- `template/` : templates HTML (fragments)

## Auteur
- Giampaolo Tranchida
