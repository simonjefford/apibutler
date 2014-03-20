/* global App, Rickshaw */

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
