angular.module('tasks', []).config(function($interpolateProvider) {
  $interpolateProvider.startSymbol('{[');
  $interpolateProvider.endSymbol(']}');
});

function TaskCtrl($scope, $http) {
    $http.get('/tasks/list').success(function(data) {
        $scope.tasks = data
    })
}