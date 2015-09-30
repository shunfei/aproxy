# aproxy

`aproxy` is a reverse proxy that includes authentication. It is designed to protect the resources that you want to expose, but only allow some one has you permission to access.

## Screenshot

**Backend config**:

![](doc/img/backend.png)

**Role List**:

![](doc/img/role.png)

**Authority config**:

![](doc/img/authority.png)


## Install

### Install from source

```
cd $GOPATH/src
git clone https://github.com/shunfei/aproxy.git
cd aproxy
sh ./install.sh
```

### Install from tarball

Go to [releases](https://github.com/shunfei/aproxy/releases) page download the tar file.

```
tar xzvf aproxy-v0.1-xxxx-xxx-xx.tar.gz
cd aproxy-v0.1-xxxx-xxx-xx
cp conf/aproxy.toml.example conf/aproxy.toml
```

## Run

Before running, your need set up [MongoDB](http://docs.mongodb.org/manual/installation/) and [Redis](http://redis.io/download#installation) (MongoDB for config storage, Redis for session storage),
and change the config in `conf/aproxy.toml`.

```
./bin/aproxy -c conf/aproxy.toml
```

By now there is no users in the database, so let me add a user:

```
./bin/adduser -c conf/aproxy.toml -action adduser -email yourname@gmail.com -pwd passwordxxx
```

And the user added above do not have admin permission, so let me set it to admin.

```
./bin/adduser -c conf/aproxy.toml -action setadmin -email yourname@gmail.com -adminlevel 99
```

And now you can visit `http://127.0.0.1:8098/-_-aproxy-_-/` and config your aproxy.

## Config

 `conf/aproxy.toml` 

## Nginx Config Example

Assuming that the resources required authorized all are the domain of `pri.domain.com`'s subdomain,
Aproxy nginx server configuration should look like:

```
server {
  listen 80;
  server_name pri.domain.com *.pri.domain.com;

  location / {
    include proxy.conf;
    # pass to aproxy
    proxy_pass http://127.0.0.1:8098;
  }

}
```

And then set the WildCard DNS Record `*.pri.domain.com` to this nginx server.

Assume that we have the following domain:

- pri.domain.com
- hadoop.pri.domain.com
- druid.pri.domain.com
- aerospike.pri.domain.com

Then we can set the login domain to `pri.domain.com`, to ensure that the sub-domain of `pri.domain.com` ( for example `hadoop.pri.domain.com`) can get the session cookies after login.    
So we change `conf/aproxy.toml` to set the domain:

```
loginHost = "http://pri.domain.com"
[session]
domain = "pri.domain.com"
```

## Integration with your company's account system

Aproxy's authority is base on email, so if your company's account system has email field, can be integration.    
To integration with aproxy, just need implement the interface of `aproxy/module/auth/UserStorager`.

```
type UserStorager interface {
    Login(email, pwd string) (*User, error)
    GetByEmail(email string) (*User, error)
    GetAll() ([]User, error)
    // add new user.
    // user.Pwd field has encrypted.
    Insert(user User) error
    Update(id string, user User) error
}
```

If you don't need manage the user in aproxy, you can just implement the `Login(email, pwd string) (*User, error)` func. 

After implement the `aproxy/module/auth/UserStorager` interface, we need change the code in `aproxy/bin/main.go`:

```
//file: aproxy/bin/main.go

delete this line:
//auth.SetUserStorageToMongo()

add this code, to register your own UserStorager to aproxy
auth.SetUserStorage(&yourUserStorage{})
```
