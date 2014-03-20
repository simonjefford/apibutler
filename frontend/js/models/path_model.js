/* global App, ajaxWithWrapperObject, ajax */

App.Path = Ember.Object.extend({
    objectForSaving: function() {
        return {
            fragment: this.get('fragment'),
            limit: parseInt(this.get('limit'), 10),
            seconds:  parseInt(this.get('seconds'), 10)
        };
    }.property('fragment', 'limit', 'seconds'),

    save: function() {
        return ajax.request('paths', {
            data: JSON.stringify(this.get('objectForSaving')),
            type: 'POST',
            dataType: 'json'
        });
    }
});

App.Path.reopenClass({
    findAll: function() {
        return ajaxWithWrapperObject('/paths', App.Path);
    }
});
