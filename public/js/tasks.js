angular.module('tasks', ['hammerServices', 'services.notifications'])
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

function CreateTaskCtrl($scope, $location, Task) {
    $scope.save = function() {
        Task.save({}, $scope.task, function(task) {
            $location.path('/');
        });
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
            notifications.send(err.data);
            $location.path('/');
        }
    );

}