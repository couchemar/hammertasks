angular.module('app', ['tasks', 'services.notifications'])
    .config(function($interpolateProvider) {
        $interpolateProvider.startSymbol('{[');
        $interpolateProvider.endSymbol(']}');
    })
    .controller(
        'AppCtrl',
        ['$scope', 'notifications',
         function($scope, notifications) {
             $scope.notifications = notifications;
             $scope.removeNotification = function(notification) {
                 notifications.remove(notification);
             };
         }
        ]);

