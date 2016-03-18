# Swan

Swan brings the power of migrations to your terminal. Swan is inspired by [Goose](https://bitbucket.org/liamstask/goose).

## Motivation

I was enjoying the simplicity and abstraction of database migrations. I was feeling the pain of distributing one off commands to update local and production environments. For local environments, I used a script to create an environment from scratch. This works well for initial creation and updates since the environment can be blown away and recreated. For production environments, I was calculating the necessary commands to update services. It really bugged me that the workflows for updating local and production environments are different. Ideally, my local and production workflows should be the same.

It hit me that the workflow for applying database updates isn't that much different than service updates. Databases are initialized once and updated every time after that. It's much easier to apply specific updates to a live database than recreating it. Does this sound familiar to deploying service updates?

Swan is a tool to package and distribute terminal commands like database migrations. Swan embraces the fact that environments are updated much more than they are created and enables you to easily distribute commands to update, or migrate, environments.

## Usage

Swan needs two things to operate: a migration directory and a last migration file. The migration directory contains executable migrations. The last migration file holds the name of the last migration. By default, the migrations directory is "migrations" and the last migration file is ".swan".

Swan currently has two commands: run and create. Run uses the migration directory and last migration file to run all migrations after the migration named in the last migration file. Create uses the migration directory and a provided name and extension to create a new file prefixed with the current timestamp.

For more usage, run `swan --help`.