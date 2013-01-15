function Task(data) {
    this.id = data.Id;
    this.summary = ko.observable(data.Summary);
    this.description = ko.observable(data.Description);
}

function TaskListViewModel() {
    var self = this;
    self.tasks = ko.observableArray([]);

    // init.
    $.getJSON("/tasks/list", function(allData) {
        var tasks = $.map(allData, function(item) {
            return new Task(item);
        });
        self.tasks(tasks);
    });
}

$(function() {
    ko.applyBindings(new TaskListViewModel());
});