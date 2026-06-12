# 🤝 Hozzájárulás a Panaszláda projekthez

Köszi, hogy érdeklődsz a Panaszláda iránt! 🚀

Ez egy nyílt forráskódú, közösségi alapú állampolgári bejelentő platform, amelynek célja egy átláthatóbb és használhatóbb hibabejelentő rendszer létrehozása.

A projekt jelenleg aktív MVP fázisban van, ezért minden segítség számít - főleg frontend és UX fejlesztés terén.

---

# 🌍 Gyors áttekintés

👉 Live demo: https://panaszlada.onrender.com

A rendszer jelenleg:

- működő backend (Go)
- működő API
- térképes bejelentés (MapLibre GL)
- alap UI

A fókusz most: **felhasználói élmény és frontend finomítás**

---

# 🚀 Hogyan tudsz csatlakozni?

## 1. Fork & clone

```bash
git clone https://github.com/istvzsig/panaszlada.git
cd panaszlada
```

## 2. Backend indítás

```bash
chmod +x ./start-dev.sh
```

A szerver:

```text
http://localhost:8080
```

Postgres:

```bash
createdb panaszlada
psql -d panaszlada -f ./db/migrate.sql
```

Data seed:

```bash
psql -d panaszlada -v count=100 -f ./db/seed_data.sql
```

Seed torles:

```bash
psql -d panaszlada -f ./db/delete_seed.sql
```

## 3. Frontend

A frontend statikus:

```text
/frontend/index.html
```

# ⚙️ Környezeti változók

Hozz létre egy .env.local fájlt:

```text
DATABASE_URL=postgresql://supabase_admin:postgres@127.0.0.1:54322/postgres (opcionalis)
PORT=8080
API_KEY=kgP4TY...
```

vagy:

```text
touch generate_env.sh
```

```bash
#!/usr/bin/env bash
# require a key name argument; generate a strong API key and insert/update it in a .env file
# Usage: bash ./gen_and_update_env.sh KEY_NAME [ENV_FILE]
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "Usage: $0 KEY_NAME [ENV_FILE]" >&2
  exit 1
fi

KEY_NAME="$1"
ENV_FILE="${2:-.env}"
KEY_BYTES=32   # 32 bytes = 256 bits

# ensure env file exists
mkdir -p "$(dirname "$ENV_FILE")" 2>/dev/null || true
if [[ ! -f "$ENV_FILE" ]]; then touch "$ENV_FILE"; fi

# generate a strong url-safe base64 key (no padding)
API_KEY=$(openssl rand -base64 "$KEY_BYTES" | tr '+/' '-_' | tr -d '=')

# backup existing env file
# cp -- "$ENV_FILE" "${ENV_FILE}.bak.$(date +%s)"

# update existing key or append if missing, preserving file format and other keys
awk -v k="$KEY_NAME" -v v="$API_KEY" -F= '
  BEGIN{updated=0}
  /^[[:space:]]*#/ { print; next }
  /^[[:space:]]*$/ { print; next }
  {
    name=$1
    gsub(/^[ \t]+|[ \t]+$/, "", name)
    if (name==k) {
      print k"="v
      updated=1
    } else {
      print
    }
  }
  END{ if (!updated) print k"="v }
' "$ENV_FILE" > "${ENV_FILE}.tmp" && mv "${ENV_FILE}.tmp" "$ENV_FILE"

echo "Updated $ENV_FILE with $KEY_NAME"
echo "$KEY_NAME=$API_KEY"

```

# 🎯 Miben tudsz segíteni?

🎨 Frontend / UI

- térkép UI finomítása
- marker + popup design javítása
- reszponzív mobil nézet
- form UX fejlesztés
- dark mode / layout polishing

# 🗺️ Map funkciók

- marker clustering (sok report esetén)
- térkép zoom / UX finomítás
- click interaction javítása
- popup usability fejlesztés

# 🧠 Feature fejlesztés

- kategória rendszer bővítése
- keresés és szűrés
- képfeltöltés
- statisztika / dashboard

# 🐛 Hibajavítás

- UI bugok
- API edge case-ek
- frontend/backend mismatch fixek

# 🧪 Jó PR szabályok

Kérlek tartsd be:

- kis, fókuszált PR-ek
- egy PR = egy dolog
- ne keverj UI + backend változtatást
- ha nagyobb változás → nyiss issue-t előtte

# 🟢 “Good First Issue” stratégia

Ha nem tudod hol kezdj:

👉 nézd a GitHub Issues részt
👉 keresd a good first issue címkéjű feladatokat

Ha nincs elég issue:
👉 nyiss egyet egy ötlettel vagy buggal

# 💬 Kommunikáció

- GitHub Issue
- vagy komment az adott issue alatt

# 🎯 Mi a cél?

Nem egy tökéletes startupot építünk.

Hanem egy:

> egyszerű, gyors, használható állampolgári bejelentő rendszert

# ❤️ Fontos

Ez egy élő, fejlődő projekt.

Minden hozzájárulás számít - legyen az:

- egy bugfix
- egy UI javítás
- vagy egy új ötlet

Köszi, hogy segítesz! 🚀
