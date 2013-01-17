angular.module('tasks', ['hammerServices'])
    .config(function($interpolateProvider) {
        $interpolateProvider.startSymbol('{[');
        $interpolateProvider.endSymbol(']}');
    })
    .config(function($routeProvider) {
        $routeProvider
            .when('/', {controller:ListCtrl, templateUrl:'tasks/list'})
            .when('/new', {controller:CreateTaskCtrl, templateUrl:'tasks/detail'})
            .otherwise({redirectTo:'/'});
    });

function ListCtrl($scope, Task) {
    $scope.tasks = Task.query();
}

function CreateTaskCtrl($scope, $location, Task) {
    $scope.save = function() {
        Task.save({}, $scope.task, function(task) {
            $location.path('/');
        })
    };
}