angular.module('services.notifications', ['ngResource'])
    .factory('notifications', function() {
        var notifications=[];
        return {
            send: function(notification) {
                notifications.push(notification);
            },
            sendSuccess: function(message) {
                notifications.push({'message': message,
                                    'type': 'success'});
            },
            sendError: function(message) {
                notifications.push({'message': message,
                                    'type': 'error'});
            },
            get: function() {
                return notifications;
            },
            remove: function(notification) {
                var idx = notifications.indexOf(notification);
                if (idx > -1) {
                    notifications.splice(idx, 1);
                }
            }
        };
    });