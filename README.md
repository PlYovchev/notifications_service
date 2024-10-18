# Notifications service assignment

## Description
This repository implements a notification service which accepts notification objects over HTTP REST and pushes them to different notification channels. Currently the supported channels are Email and Slack but the implementation allows easy extension for additional notification channels like SMS, etc. 

## Architecture

![SumUpNotificationService_v3](https://github.com/user-attachments/assets/e49c121c-ca39-4a3b-8c66-e66154d16440)

https://drive.google.com/file/d/1FKyFudmjgg_3aQjebkB7i4CTZhbPkEqL/view?usp=sharing

The project is built to run in a docker environment. It consists of the following services:
1. Nginx - reverse proxy service which is reponsible for hiding the internal API and controlling which one should be exposed publicly. Other reponsibilities for the nginx are loadbalancing and rate-limiting. The config of nginx is in /nginx folder.
2. Postgres - the persistence layer is a Postgres database.
3. Notification service - a Golang app which exposes an HTTP endpoint for pushing notifications, handles persistence of the notifications in the persistence layer, and it is responsible for scheduling the actual sending of those notifications over the supported channels.

### Notification service app
The notification service app is written in Golang. The following libraries are used:
1. Routing - [GIN](https://github.com/gin-gonic/gin)
2. Logging - [zerolog](https://github.com/rs/zerolog)
3. Database connection & ORM - [GORM](https://github.com/go-gorm/gorm)

#### The service works with these main entities:
1. Handlers - responsible for the application level logic - handling the HTTP request, validating and transforming the input and transfering it to a service entity;
2. Services - responsible for the business level logic;
3. Repositories - resposible for the persistence level logic;
4. Notifiers - reponsible for connecting and performing the actual send notification logic to 3rd party providers and services.

#### Exposed APIs:
1. **POST /public-api/v1/notifications/push-notification** - accepts a JSON NotificationInput object. Responsible for submitting a notification to be sent over the delivery channels specified in the input;
    - example usage (the snippet direct the request to the NGINX and should be executed outside of the docker env):
    ``` 
    curl -d '{ "key":"payment-cancelled","message":"Payment has failed", "deliveryChannels": ["Email", "Slack"] }' -X POST localhost:3000/v1/notifications/push-notification
    ```
2. **GET /status** - internal API which checks if the service is healthy;

#### Implementation behavior:
The behavior of the notification service app is depicted on the diagram above. The key elements are:
1. Once a notification input is pushed to the '/notifications/push-notifications' endpoint, the notification input is transformed into separate notification objects. The transformation logic uses the *notificationInput.deliveryChannels* property to determine how many notifications should be created - one for each delivery channel;
2. After the internal notification objects are created, they are persisted with status **PENDING** in the database and the polling notification service object is notified that new notifications have been received;
3. The observer/polling mechanism of the notification service is started with the starting of the app. It is responsible for processing any pending notifications that are stored in the database. It performs a polling logic over a specific period of time for any pending notifications, and it also allows to be forcefully awaken using **notificationService#OnNotificationsReceived(notificationIds)** to process and prioritize any newly arrived notifications.
4. When the notifications are processed, in case of error or missing confirmation that a specific notifier successfully sent the notication over a channel, the processing for those failed notifications is retried in total of 3 times;
5. Upon completion of sending of the notifications or exhausting the retry count, the notifications are saved in the database with updated status, respectively 'completed' and 'failed'.
6. The notification status 'completed' and 'failed' are considered terminal at the moment.

## Deployment
The configuration in the docker-compose.yaml deploys 4 services:
1. **notification_service** - this service is the Golang application. It is configured to be deployed with 3 replicas and health checks.
2. **nginx** - the reverse proxy service is configured to use the 3 replicas for loadbalancing using the least connections strategy. It is also performing rate-limiting of 50 request per second per IP.
3. **postgres** - the database service.
4. **autoheal** - this service is responsible to redeploy the **notification_service** if the service is deemed unhealthy.

## Running the project

#### Pre-requisites: **Docker**, **Docker Compose**, **Make**

1. ```make setup``` - builds the docker images.   
2. ```make start``` - the **start** target calls ``docker-compose up -d`` so docker compose has to be installed in advance.
3. ```make clean``` - stops containers and removes containers, networks, volumes, and images.

#### Configuring the **notifiers**.
The notifiers use properties which are sourced from **/resources/config/application.*.yml**. When running this setup with ``make start``, use **/resources/config/application.docker.yml**.
The following properties have to be set:
1. **EmailNotifier** required data:
    - **from** - the email address from which the email notifications should be sent;
    - **password** - password for the *from* email address (for gmail.com this password should be generated by 'App password' functionality)
    - **recepients** - the email addresses of the receivers for the notifications;
    - **smtpHost** & **smtpPort** - host and port of the SMTP server;
2. **SlackNotifier** required data:
    - **webhookUrl** - valid webhook url generated by the 'https://api.slack.com/apps/' for the specific channel in Slack that should receive the notifications;

## TODO
1. Add unit tests as the key components of the notification service app are not covered with unit tests yet;
2. Add Kubernetes deployment scripts & configuration;
3. Add Github actions/prehooks for running lint&tests on commits;
4. Implement some authentication mechanism for the notification app - example apiKey, etc;
5. Improve the error handling;
