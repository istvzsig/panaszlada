# 🔐 Biztonsági irányelvek (Security Policy)

## 📌 Áttekintés

A Panaszláda rendszer közérdekű bejelentéseket kezel, amelyek tartalmazhatnak érzékeny információkat.

A biztonság és az anonimitás kiemelten fontos tervezési elv.

---

## 🔑 Titkok kezelése (Secrets Management)

- Soha ne commitolj `.env` vagy egyéb credential fájlokat
- Minden érzékeny adat környezeti változókból (`environment variables`) érkezik
- API kulcsokat és titkokat rendszeresen forgatni kell (rotate)
- Erős, véletlenszerű értékeket kell használni (pl. API_KEY, SESSION_SECRET)

---

## 🛡️ Adatvédelem

- Nem tárolunk szükségtelen személyes adatokat
- Az adatminimalizálás elvét követjük
- Ha érzékeny adat kerül tárolásra, azt titkosítani kell
- Adatok továbbítása kizárólag HTTPS-en keresztül történik

---

## 🔐 Hitelesítés és jogosultságkezelés

- Csak HTTPS kapcsolat használható éles környezetben
- JWT alapú autentikáció használata (lejárati idővel)
- Cookie-k esetén: `httpOnly`, `Secure`, `SameSite` beállítások
- Jogosultságok szigorú ellenőrzése minden admin endpointon
- CORS szabályok helyes konfigurálása

---

## 🕵️ Anonimitás

- IP címek nem kerülnek nyers formában tárolásra
- UUID alapú azonosítók használata
- Feltöltött fájlok metaadatainak eltávolítása
- Audit logok nem tartalmazhatnak személyes adatokat
- Lehetőség szerint mezőszintű titkosítás

---

## 🚦 Rate limiting és védelem

- API limit: pl. 100 kérés / 15 perc / IP
- Érzékeny endpointok védelme (pl. CAPTCHA)
- Reverse proxy / WAF használata éles környezetben

---

## 💻 Kódbiztonság

- Függőségek rendszeres frissítése
- Vulnerability scan használata
- Csak ellenőrzött csomagok használata
- Input validáció minden endpointon
- Output sanitization XSS ellen

---

## 📊 Naplózás és monitorozás

- Biztonsági események naplózása
- Személyes adatokat nem tartalmazhat log
- Éles rendszerben csak minimális log szint
- Rendszeres log ellenőrzés

---

## 🚀 Deployment biztonság

- HTTPS kötelező
- Biztonságos HTTP headerek használata
- Környezetenként külön konfiguráció
- Adatbázis SSL kapcsolat használata
- Rendszeres patch-elés

---

## 🧠 API biztonság

- Server-side input validáció kötelező
- XSS elleni védelem
- CSRF védelem state-changing műveleteknél
- API verziózás javasolt
- Publikus dokumentáció a használathoz

---

## 🗄️ Adatbázis biztonság

- Erős jelszavak használata
- Least privilege elv
- SSL kapcsolat adatbázis felé
- Backup-ok titkosítása
- Hozzáférések rendszeres auditálása

---

## 🚨 Hibabejelentés (Security Reporting)

Ha biztonsági hibát találsz:

👉 kérlek ne nyiss public issue-t  
👉 írj emailt: **security@example.com**

A hibákat felelősségteljesen kezeljük és javítás után kommunikáljuk.

---

## 📋 Fejlesztői ellenőrző lista

- [ ] Nincs secret commitolva
- [ ] Minden input validálva
- [ ] Nincs PII felesleges gyűjtése
- [ ] Hibák nem szivárogtatnak érzékeny adatot
- [ ] Függőségek ellenőrizve
- [ ] OWASP alapelvek követve

---

## 📌 Megjegyzés

Ez a projekt aktív fejlesztés alatt áll, a biztonsági gyakorlatok a jövőben finomodhatnak.
