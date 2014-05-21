var ConfigTextField = Ember.TextField.extend({
    fieldName: null,

    config: null,

    valueChange: function() {
        var fieldName = this.get('fieldName');
        var config = this.get('config');
        if (fieldName && config) {
            config[fieldName] = this.get('value');
        }
    }.observes('value'),

    loadValue: function() {
        var fieldName = this.get('fieldName');
        var config = this.get('config');
        if (fieldName && config) {
            this.set('value', config[fieldName]);
        }
    }.on('didInsertElement')
});

export default ConfigTextField;
