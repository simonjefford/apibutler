/* global Rickshaw */

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
