# Mon Caveau

Mon caveau est un site qui permet de gérer sa cave (ou "caveau"). <br/>
Ce projet n'est plus en développement et a seulement été créé pour satisfaire un besoin d'un proche. <br/>
Si vous souhaitez tester avec des données : <br/>
- Accédez à `{Nom de domaine [www.moncaveau.com] ou votre installation particulière)/debug/fakeit`.
- Récupérez la clé générée et connectez-vous (`/v1/login`).

## TODO

- [X] Database system with MySql
    - [X] Env usage for connection
    - [X] Auto migrations
    - [X] Entities dynamic sql creation
        - [X] Avoid inserting an id when creating a new entity
- [X] Server
    - [X] Quick starter with config (env)
    - [X] Logger system
    - [X] Crash handler attached
    - [X] Quit handler to clean everything
    - [X] Allow TLS config (aka https)
- [X] Authentification 
- [ ] Frontend
    - [X] Debug system with direct file access
    - [X] Release mode with no dependency (use embed.FS to keep files in memory (faster))
    - [X] Create the front page
    - [X] Create the login page 
    - [ ] Create the register page
    - [X] Wine dashboard/dashboard page
    - [X] Wine dashboard/collection page
    - [X] Wine dashboard/statistics page
    - [ ] Wine dashboard/account page
    - [ ] All the pages above should be nice to use either on phone and computer
- [ ] Login / Register
    - [X] Database system
    - [ ] Allow no login / register but create a temporary user with a infinite cookie
    - [X] Modern login system (no username / password)
    - [ ] Allow to add email to the account and be able to log in with it
    - [ ] Allow to add a password and be able to login with it 
    - [X] Login with only Account Key (generated unique)
    - [X] Password encryption with salt (+ peper ?)
    - [X] Basic cookie system with expiration
    - [X] Cookie system with "remember me" thing
    - [ ] Email verification
    - [ ] Allow to force / unforce some fields for the login (password / email / account key) (at least 1 and password / email can't be alone) 
- [ ] Wine collection system
    - [X] Basic database system 
    - [X] Wines linked to a user 
    - [X] We can add a new region
    - [X] We can add a new type of wine
    - [X] We can add a new domain
    - [X] We can add a new bottle size
    - [ ] Allow all the adds above with verification by admin
    - [X] Create statistics from the wines
    - [ ] Allow quick modification of quantity
    - [ ] Allow reusage of images of the same wine name for an other user
    - [ ] Allow auto completion when typing the wine name (offer the posibility to use already created wine name and complete all the other fields with the correct data but let the user obviously change everything if he wants)
- [ ] Wine rating system
    - [ ] Be able to rate a wine very fast and easly
- [ ] Wine label scan
    - [ ] Train AI ?
- [ ] Email system
    - [ ] Notification system (extra)
    - [ ] Email verification (login)
    - [ ] Setup
        - [ ] Setup a Mail Transfer Agent
        - [ ] Setup a Mail Delivery Agent
        - [ ] Setup a SPF dns setting
        - [ ] Setup a DKIM dns setting
        - [ ] Setup a DMARC dns setting
