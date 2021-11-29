# Github action watcher

Basic objective of this project is to re-run workflow if that stucks at any point of execution. this is temparary solution to keep watching github actions and trigger cancell and rerun as needed.

### Generate token

To run this application you supposed to be generate respective github token (to access private repo on behalf of you)

Use this link to generate token
[https://github.com/settings/tokens](https://github.com/settings/tokens)


### Configure

Basically we need to set the generated token to `.env` file. for that you need to copy `.env.example` file to `.env` file. and add your token and username there

```bash
cp .env.example .env
vi .env
# or you can open it in any of the text editor of your choice
# and update the values of token and username
```

Another configuration that is required is to set repo and it's max time limit to allow github actions to run.

for that you can copy `conf.example.json` to `conf.json`

```bash
cp conf.example.json conf.json
```

and make respective changes.

**Name** - `string` Repo name will be combination of owner and name if you take this repo as an example then then name field should have value as `AdiechaHK/gha-watcher`

**Max Minutes** - `int` Maximum time that we allow github actions to be run this will be in integer and that will indecats minutes to allow run and wait.


### Execute 

##### for development

```bash
go run main.go
```

##### for production

```bash
go run build
./gha-watcher
```