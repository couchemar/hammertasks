angular.module('resources.tasks', ['ngResource'])
    .factory('Task', function($resource) {
        var Task = $resource(
            'tasks/json/:id',
            {id: '@id'},
            {update: {method: 'PUT'}}
        );

        Task.prototype.update = function(success, fail) {
            return Task.update(
                {id: this.id},
                angular.extend({}, this),
                success,
                fail);
        };
        Task.prototype.remove = function(success, fail) {
            return Task.remove(
                {id: this.id},
                success,
                fail
            );
        };
        return Task;
    });