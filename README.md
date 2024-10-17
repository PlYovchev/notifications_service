# Notifications service assignment

## Description
This repo implements a notification service which accepts notification objects over HTTP REST and pushes them to different notification channels. Currently the supoorted channels are Email and Slack but the implementation allows for easy extension of additional notification channels like SMS, etc. 

## Architecture

![SumUpNotificationService](https://github.com/user-attachments/assets/8580884e-fa08-41e6-86cc-650abfa7b17a)

https://drive.google.com/file/d/1FKyFudmjgg_3aQjebkB7i4CTZhbPkEqL/view?usp=sharing

The project is built to run in a docker environment. It consists of the following services:
1. Nginx - reverse proxy service which is reponsible for hiding the internal API and controlling which one should be exposed publicly. Other reponsibilities for the nginx are loadbalancing and rate-limiting.
2. Postgres - the persistence layer is a Postgres database.
3. Notification service - a Golang app which exposes an HTTP endpoint for pushing notifications, handles persisting the notifications in the persistence layer, and it is responsible for scheduling the actual sending of those notifications over the supported channels.

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
1. Once a notification input is pushed to the '/notifications/push-notifications' endpoint, the notification input is transformed into separate notification objects. The transformation logic uses the notificationInput.deliveryChannels to determine how many notifications should be created - one for each delivery channel;
2. After the internal notification objects are created, they are persisted with status **PENDING** in the database and the polling notification service object is notified that new notifications have been received;
3. The notification service object is started with the starting of the app. It is responsible for processing any pending notifications that are stored in the database. It performs a polling logic over a specific period of time for any pending notifications, and it allows to be forcefully awaken using **notificationService#OnNotificationsReceived(notificationIds)** to process and prioritize any newly arrived notifications.
4. When the notifications are processed, in case of error or misssing confirmation when a specific notifier is attempting to send notication over a channel, the processing for those failed notifications is retried in total of 3 times;
5. Upon completion of sending of the notifications or exhausting the retry count, the notifications are saved in the database with updated status, respectively 'completed' or 'failed'.
6. The notification status 'completed' and 'failed' are considered terminal at the moment.

## Running the project

Pre-requisites: Docker, Docker Compose, Make

1. 
