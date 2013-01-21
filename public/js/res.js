angular.module('hammerServices', ['ngResource'])
    .factory('Task', function($resource) {
        var Task = $resource(
            'tasks/json/:id',
            {id: '@id'},
            {update: {method: 'PUT'}}
        );
        return Task;
    });