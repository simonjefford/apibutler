var ConfigTextField = Ember.TextField.extend({
    fieldName: null,

    config: Ember.computed.alias('parentView.config'),

    configParent: Ember.computed.alias('parentView.configParent'),

    classNameBindings: [':config_field', 'fieldName'],

    valueChange: function() {
        var fieldName = this.get('fieldName'),
            config = this.get('config'),
            configParent = this.get('configParent');
        if (fieldName && config) {
            config.set(fieldName, this.get('value'));
        }
        configParent.incrementProperty('_changes');
    }.observes('value'),

    loadValue: function() {
        var fieldName = this.get('fieldName');
        var config = this.get('config');
        if (fieldName && config) {
            this.set('value', config.get(fieldName));
        }
    }.on('didInsertElement')
});

export default ConfigTextField;
