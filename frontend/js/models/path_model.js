/* global App, DS */

App.Path = DS.Model.extend({
    fragment: DS.attr('string'),
    limit: DS.attr('number'),
    seconds: DS.attr('number')
});
