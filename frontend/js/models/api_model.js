/* global App, DS */

App.Api = DS.Model.extend({
    fragment: DS.attr('string'),
    limit: DS.attr('number'),
    seconds: DS.attr('number')
});
