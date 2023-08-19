# oj-lab-services

This is a collection of services for OJ-lab systems.

Each service is supposed to run either separately or in a batch,
so that the deployment of the online-judge system can be more flexible.

## Services

To understand the philosophy of our design,
you can simply think an online-judge system consists of three main parts:
- user
- problem
- judge

If every request is trusted
(and thanks to the usage of JWT authentication, this is now a solved problem),
each of the above parts can be considered as a single service. 

So separating them can make the build of functionality clearer than ever before.

## Model migration

There is a migration script in every service directory.

You will need to run the migration script before you can use the service.

## Development

For service development, we don't want to make it too complex.
Using VSCode on either Win/*nix System are avaliabe, try using the Makefile/Dockerfile in the repository.