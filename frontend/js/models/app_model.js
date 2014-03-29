/* global App, DS */

App.App = DS.Model.extend({
    name : DS.attr('string'),
    backendURL: DS.attr('string')
});
