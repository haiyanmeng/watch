Follow the following step to test this:

Step 1: replace the docker image path in `config.yaml` and `build.sh`;

Step 2: run `kubectl delete -f config.yaml` to delete the job if it exists;

Step 3: run `./build.sh` to build/push the docker image and start the job.

Step 4: run `kubectl logs job.batch/watch` to get the log.
