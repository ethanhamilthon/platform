## Todo
1. Add Jetstream and use it as default, where we need response (rest api GET request alternative) use just core
2. Add Env and Config modules, topics and consumer name should be passed with Env
3. Run the platform in http 80 port by default. Pipe: 1. Run Platform server, 2. Run UI, 3. Run Http 80 server, 4. Add Platform and UI as apps. If there are already domains we will go through another path. 
4. Btw dont forget to use volumes instand of binds (in platform too, we wont give accesses to user to file system. Instand of that we will create volumes, and user can connect them in their apps) 


## Somethings

### how manager works
ON START:
    - Check if 443 and 80 ports are available
    - Run UI, Controller and Adapter containers
    - Get balancer state from Controller
    - Run Controller API 
    - Run Balancer container 
    - Send states to Balancer service
    - Run http/https balancers
