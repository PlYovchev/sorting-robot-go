# sort

Steps to build and run the container used to build the app and run it:

1. docker build -t sorting-robot-app .
2. docker run -i -t -v <path-to-the-source-base-path>:/go/src/app -p 10000:10000 -p 10001:10001 sorting-robot-app
