# Mon Caveau

## TODO

- [X] Database system with MySql
    - [X] Env usage for connection
    - [X] Auto migrations
    - [X] Entities dynamic sql creation
        - [X] Avoid inserting an id when creating a new entity
- [X] Server starter + logger + crash handler + quit handler
- [ ] Authentification 
- [ ] Frontend
    - [X] Debug system with direct file access
    - [X] Release mode with no dependency (use embed.FS to keep files in memory (faster))
    - [ ] Create the front page
    - [ ] Create the login / register page
- [ ] Login / Register
    - [ ] Database system
    - [ ] Allow no login / register but create a temporary user with a infinite cookie
    - [X] Modern login system (no username / password)
    - [ ] Allow to add email to the account and be able to log in with it
    - [ ] Allow to add a password and be able to login with it 
    - [X] Login with only Account Key (generated unique)
    - [X] Password encryption with salt (+ peper ?)
    - [ ] Cookie system with "remember me" thing
    - [ ] Email verification
    - [ ] Allow to force / unforce some fields for the login (password / email / account key) (at least 1 and password / email can't be alone) 
- [ ] Email system
    - [ ] Notification system (extra)
    - [ ] Email verification (login)
    - [ ] Setup
        - [ ] Setup a Mail Transfer Agent
        - [ ] Setup a Mail Delivery Agent
        - [ ] Setup a SPF dns setting
        - [ ] Setup a DKIM dns setting
        - [ ] Setup a DMARC dns setting
        - [ ] Setup TLS encryption