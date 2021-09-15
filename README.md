# userProfile

# todo

- support writeClient and readClient for postgres
- redis Replicas read from config
- config hot reload
- complete readme
- staging and canary support :)

## Quick start & development

First, you need to create a database and user. Then export values of `DATABASE_USER` and `DATABASE_PASSWORD` to the
shell. You can do this by following commands:

```bash
cp .sample.env .env
export $(xargs <.env)
```

In this case, you need to have a username and password just like those specified in `.env` file that you can
successfully connect to DB.

Then you can build and serve project with these commands.

```bash
docker-compose up --force-recreate
make dependencies
make userProfile 
```

If there were issues on MacOS run the followings:

```
brew install librdkafk

echo 'export PATH="/opt/homebrew/opt/openssl@1.1/bin:$PATH"' >> .zprofile
export 'PKG_CONFIG_PATH="/opt/homebrew/opt/openssl@1.1/lib/pkgconfig"' >> .zprofile

```

If there were issues on make dependencies

```bash
1 - make sure your vpn is on :) 
```

## after your work :

```bash
docker-compose down --remove-orphans -v
```