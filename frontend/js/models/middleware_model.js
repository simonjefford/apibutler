App.Middleware = DS.Model.extend({
    friendlyName: DS.attr('string'),
    configItems: DS.attr(),
    name: DS.attr('string')
});
