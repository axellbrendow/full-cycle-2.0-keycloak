# full-cycle-2.0-keycloak

Files I produced during the Keycloak / OAuth2 / OpenID Connect classes of my [Microservices Full Cycle 3.0 course](https://drive.google.com/file/d/1bJnFxQPKgSsI30sCvW-KzYK4V5JWzgSs/view?usp=share_link).

## Running Keycloak container

```sh
docker run -p 8080:8080 -e KEYCLOAK_USER=admin -e KEYCLOAK_PASSWORD=admin quay.io/keycloak/keycloak
```

## Creating a Realm

Realm is a concept in Keycloak that refers to an object managing a set of users along with their credentials, roles and groups. A Realm works like a tenant. Multiple clients can be interested in accessing users data from a Realm.

Access http://localhost:8080/auth/admin/ on your browser and login with admin/admin.

The Master Realm is used to manage Keycloak itself. We will create a new Realm from it.

![Shows the Master Realm screen and the add realm button](./imgs/master-realm.png)

Use myrealm as the name and create.

## Creating an User

Now, create an user inside myrealm.

![Shows how to create the user inputting the username, first name and last name](./imgs/create-user.png)

And set its credentials.

![Shows how to set a password for the user](./imgs/set-credentials.png)

## Creating a Client

Create a client for our Go application.

![Shows how to create a client that will access that user info](./imgs/create-client.png)

After creating, change the client access type to confidential and save.

![Shows how to change the client access type to confidential](./imgs/change-access-type.png)

Get the client secret and replace in `goclient/main.go`.

![Shows the Credentials tab where you can find the client secret](./imgs/get-client-secret.png)

![Shows where to change the client secret in the Go application code](./imgs/client-secret-in-goapp.png)

You also need to get the issuer URL from the Realm Settings and replace in `goclient/main.go`.

![Shows the realm settings with the link to the realm endpoints](./imgs/realm-settings-endpoints.png)

![Shows the issuer link inside the json with the endpoints](./imgs/realm-endpoints-json.png)

![Shows where to change the realm issuer URL in the Go application code](./imgs/realm-issuer-in-goapp.png)

## Running the client

```sh
go mod tidy
go run goclient/main.go
```

Access http://localhost:8081/ on your browser and login in with the same user you created in the realm.

![Shows login callback screen after a successful login](./imgs/access-token.png)

You can go to https://jwt.io to see the access token payload.

![Shows the access token payload](./imgs/access-token-payload.png)

As you can see, this step is the **authorization** step. We can now request the id token as we have access to the openid scope. With the id token, we are **authenticated**.

![Shows the id token](./imgs/id-token.png)

![Shows the id token payload](./imgs/id-token-payload.png)
