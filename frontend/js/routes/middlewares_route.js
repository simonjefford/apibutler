/* global App */

App.MiddlewaresRoute = Ember.Route.extend({
    model: function() {
        return [
            Ember.Object.create({
                name: 'Authorisation',
                count: 10
            }),
            Ember.Object.create({
                name: 'Throttling',
                count: 50,
                configSettings: [
                    'timeInterval',
                    'callCount'
                ]
            }),
        ];
    }
});
