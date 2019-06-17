echo "runtime: go111
instance_class: F1
env_variables:
  PORT: '$PORT'
  SIGNING_SECRET: '$SIGNING_SECRET'
  BOT_TOKEN: '$BOT_TOKEN'
  CACHE_DEFAULT_EXPIRATION: '$CACHE_DEFAULT_EXPIRATION'
  CACHE_CLEANUP_INTERVAL: '$CACHE_CLEANUP_INTERVAL'
  BOT_ID: '$BOT_ID'" > app.yaml

echo "successfully created app.yaml, deploying app"

gcloud app deploy app.yaml

echo "deploy finished, deleting app.yaml"

rm app.yaml
