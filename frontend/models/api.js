import ModelMixin from 'apibutler/mixins/model';

var Api = DS.Model.extend(ModelMixin, {
    path: DS.attr('string'),
    needsAuth: DS.attr('boolean'),
    app: DS.belongsTo('app'),
    stack: DS.belongsTo('stack'),

    valid: function() {
        return !Ember.isEmpty(this.get('app')) &&
            !Ember.isEmpty(this.get('path')) &&
            !Ember.isEmpty(this.get('stack'));
    }.property('app', 'path', 'stack')
});

export default Api;
