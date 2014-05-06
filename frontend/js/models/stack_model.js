App.Stack = DS.Model.extend({
    name: DS.attr('string'),
    middlewares: DS.attr(),
    configs: DS.attr()
});
