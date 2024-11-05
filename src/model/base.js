class Base {
    constructor(tableName, primaryKey, fields, fieldTypes = {}) {
        this.tableName = tableName;
        this.primaryKey = primaryKey;
        this.fields = fields;
        this.fieldsType = this.setFieldTypes(fieldTypes);
        this.formData = {};
    }

    setFieldTypes(fieldTypes) {
        return this.fields.reduce((types, field) => {
            types[field] = fieldTypes[field] || { type: 'string', value: '' };
            return types;
        }, {});
    }

    insert() {
        const columns = this.fields.filter(f => f !== this.primaryKey).join(', ');
        const values = this.fields.filter(f => f !== this.primaryKey)
            .map(f => this.escapeValue(this.formData[f], f)).join(', ');
        return `INSERT INTO ${this.tableName} (${columns}) VALUES (${values})`;
    }

    update() {
        const setClause = this.fields.filter(f => f !== this.primaryKey)
            .map(f => `${f} = ${this.escapeValue(this.formData[f], f)}`).join(', ');
        const whereClause = `${this.primaryKey} = ${this.escapeValue(this.formData[this.primaryKey], this.primaryKey)}`;
        return `UPDATE ${this.tableName} SET ${setClause} WHERE ${whereClause}`;
    }

    delete() {
        return `DELETE FROM ${this.tableName} WHERE ${this.primaryKey} = ${this.escapeValue(this.formData[this.primaryKey], this.primaryKey)}`;
    }

    selects() {
        return `SELECT ${this.fields.join(', ')} FROM ${this.tableName}`;
    }

    escapeValue(value, field) {
        const type = this.fieldsType[field]?.type;
        if (type === 'string' || type === 'enum' || type === 'datetime') {
            return `'${String(value).replace(/'/g, "''")}'`;
        }
        if (type === 'int') {
            return parseInt(value, 10) || 0; // Default to 0 if NaN
        }
        return value;
    }
}

export default Base;
