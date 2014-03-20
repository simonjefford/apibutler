/* global Spinner, Rickshaw */
var ajax = ic.ajax;

var App = Ember.Application.create({
    LOG_TRANSITIONS: true
});

App.Router.map(function() {
    this.resource('paths', function() {
        this.route('new');
    });
    this.resource('applications', function() {
        this.route('new');
    });
});

var ajaxWithWrapperObject = function(path, klass) {
    return ajax.request(path).then(function(array) {
        return array.map(function(item) {
            return klass.create(item);
        });
    });
};

App.PathsRoute = Ember.Route.extend({
    model: function() {
        return App.Path.findAll();
    }
});

App.ApplicationsRoute = Ember.Route.extend({
    model: function() {
        return new Ember.RSVP.Promise(function(resolve) {
            window.setTimeout(function(){
                var random = new Rickshaw.Fixtures.RandomData(2000);
                var data = [[]];
                for (var i=0; i<2000; i++) {
                    random.addData(data);
                }
                resolve(data[0]);
            }, 1);
        });
    }
});

App.ApplicationsController = Ember.ArrayController.extend({
    renderer: 'line',

    renderers: ['area', 'line', 'bar', 'scatterplot']
});

App.Path = Ember.Object.extend({
    objectForSaving: function() {
        return {
            fragment: this.get('fragment'),
            limit: parseInt(this.get('limit'), 10),
            seconds:  parseInt(this.get('seconds'), 10)
        };
    }.property('fragment', 'limit', 'seconds'),

    save: function() {
        return ajax.request('paths', {
            data: JSON.stringify(this.get('objectForSaving')),
            type: 'POST',
            dataType: 'json'
        });
    }
});

App.Path.reopenClass({
    findAll: function() {
        return ajaxWithWrapperObject('/paths', App.Path);
    }
});

App.NavbarLinkComponent = Ember.Component.extend({
    tagName: '',
});

App.LoadingSpinnerComponent = Ember.Component.extend({
    spinner: undefined,
    lines: 13,
    length: 20,
    width: 10,
    radius: 30,

    showSpinner: function() {
        var target = this.get('element');
        this.spinner = new Spinner({
            lines: this.get('lines'),
            length: this.get('length'),
            width: this.get('width'),
            radius: this.get('radius')
        });
        this.spinner.spin(target);
    }.on('didInsertElement'),

    teardown: function() {
        if (this.spinner) {
            this.spinner.stop();
        }
    }.on('willDestroyElement')
});

App.DataChartComponent = Ember.Component.extend({
    data: [],

    renderer: 'bar',

    width: 750,

    height: 500,

    color: 'steelBlue',

    xAxisShowsTime: true,

    xAxisTimeUnit: 'day',

    xAxis: function(graph) {
        if (this.get('xAxisShowsTime')) {
            return new Rickshaw.Graph.Axis.Time({
                graph: graph
            });
        } else {
            return new Rickshaw.Graph.Axis.X({
                graph: graph
            });
        }
    },

    showGraph: function() {
        var element = this.get('element');
        element.innerHTML = '';
        var graph = new Rickshaw.Graph({
            element: element,
            width: this.get('width'),
            height: this.get('height'),
            series: [{data: this.data, color: this.get('color')}],
            renderer: this.get('renderer')
        });

        var xAxis = this.xAxis(graph);

        xAxis.render();

        graph.render();
    }.on('didInsertElement').observes('data', 'renderer')
});

App.PathsNewRoute = Ember.Route.extend({
    model: function() {
        return App.Path.create();
    },

    actions: {
        save: function(model) {
            var self = this;
            console.log('Now saving %o', JSON.stringify(model));
            model.save().then(function(result) {
                self.controllerFor('paths').addObject(model);
                self.transitionTo('paths.index');
                return result;
            }, function(arg) {
                console.log(arg);
            });
        }
    }
});
