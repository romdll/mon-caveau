import requests

testAccountKey = "418860-750366-546845" # Mettre un compte de test
baseUrl = "http://localhost:80" # Mettre l'url de base
with requests.Session() as session:

    res = session.post(baseUrl+"/api/login", json={"account_key":testAccountKey})
    print(res)

    if res.status_code != 200:
        print("Bad login")
        exit(1)

    print(session.cookies)

    test = session.get(baseUrl+"/api/wines/basic")
    print(test)
    print(test.text)