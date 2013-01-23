angular.module('tasks', ['resources.tasks', 'services.notifications'])
    .config(function($routeProvider) {
        $routeProvider
            .when('/', {controller:ListCtrl,
                        templateUrl:'tasks/list'})
            .when('/new', {controller:CreateTaskCtrl,
                           templateUrl:'tasks/detail'})
            .when('/edit/:taskId', {controller:EditTaskCtrl,
                                    templateUrl:'tasks/detail'})
            .otherwise({redirectTo:'/'});
    });

function ListCtrl($scope, Task) {
    $scope.tasks = Task.query();
    $scope.orderProp = 'id';
}

function CreateTaskCtrl($scope, $location, Task, notifications) {
    $scope.save = function() {
        Task.save({},
                  $scope.task,
                  function(task) {
                      notifications.sendSuccess('Successfully created');
                      $location.path('/');
                  },
                  function(err) {
                      notifications.sendError('Could not create');
                      $location.path('/');
                  }
                 );
    };
}

function EditTaskCtrl($scope, $location, $routeParams,
                      Task, notifications) {
    Task.get(
        {id: $routeParams.taskId},
        function(task) {
            $scope.task = new Task(task);
        },
        function(err) {
            notifications.sendError(err.data['message']);
            $location.path('/');
        }
    );
    $scope.save = function() {
        $scope.task.update(
            function() {
                notifications.sendSuccess('Successfully saved');
                $location.path('/');
            },
            function(err) {
                notifications.sendError(!!err.data['message']?err.data['message']:'Could not save');
            });
    };
    $scope.remove = function() {
        $scope.task.remove(
            function() {
                notifications.sendSuccess('Successfully deleted');
                $location.path('/');
            },
            function(err) {
                notifications.sendError(!!err.data['message']?err.data['message']:'Could not delete');
            }
        );
    };
}