/* global App, DS */

App.Api = DS.Model.extend({
    fragment: DS.attr('string'),
    limit: DS.attr('number'),
    seconds: DS.attr('number'),
    isPrefix: DS.attr('boolean'),
    needsAuth: DS.attr('boolean'),
    app: DS.belongsTo('app')
});
