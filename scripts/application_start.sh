#!/bin/bash

#give permission for everything in the express-app directory
sudo chmod -R 777 /home/ec2-user/canvas-test

#navigate into our working directory where we have all our github files
cd /home/ec2-user/canvas-test

#docker-compose
sudo docker-compose up