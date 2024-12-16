import requests

testAccountKey = "382011-919599-481383" # Mettre un compte de test
baseUrl = "http://localhost:80" # Mettre l'url de base
with requests.Session() as session:

    res = session.post(baseUrl+"/api/login", json={"account_key":testAccountKey})
    print(res)

    if res.status_code != 200:
        print("Bad login")
        exit(1)

    print(session.cookies)

    testWine = {
		"name": "test",

		"domain": {
			"id": 1
            # "name": "test"
		},
		"region": {
            "id": 1
			# "name": "test",
			# "country": "test"
		},
		"type": {
			"id": 1
            # "name": "test"
		},
		"bottle_size": {
			"id": 1
            # "name": "test",
            # "size": 1000
		},

		"vintage": 2000,
		"quantity": 20,

		"buy_price": 1.0,
		"description": "test",
		"Image": "test"
	}
    testWineInsertion = session.post(baseUrl + "/api/wines/create", json=testWine)
    print(testWineInsertion)
    print(testWineInsertion.text)

    test = session.get(baseUrl+"/api/wines/basic")
    print(test)
    print(test.text)