var isEmpty = Ember.isEmpty, keys = Ember.keys;

App.ModelMixin = Ember.Mixin.create({
    isClean: function() {
        return isEmpty(keys(this.changedAttributes()));
    }.property().volatile()
});
