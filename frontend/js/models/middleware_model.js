App.Middleware = DS.Model.extend({
    friendlyName: DS.attr('string'),
    schema: DS.attr(),
    name: DS.attr('string')
});
