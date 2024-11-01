class base {
    tableName;
    primaryKey;
    fields;
    formData;

    constructor(tableName, primaryKey, fields) {
        this.tableName = tableName;
        this.primaryKey = primaryKey;
        this.fields = fields;
        this.formData = {};
    }

    insert() {
        const columns = this.fields.slice(1).join(', ');
        const values = this.fields.slice(1).map(field => `'${this.formData[field]}'`).join(', ');
        return `INSERT INTO ${this.tableName} (${columns}) VALUES (${values})`;
    }

    update() {
        const setClauses = this.fields.map(field => `${field} = '${this.formData[field]}'`).join(', ');
        return `UPDATE ${this.tableName} SET ${setClauses} WHERE ${this.primaryKey} = '${this.formData[this.primaryKey]}'`;
    }

    delete() {
        return `DELETE FROM ${this.tableName} WHERE ${this.primaryKey} = '${this.formData[this.primaryKey]}'`;
    }

    select_list() {
        return `SELECT ${this.fields.join(', ')} FROM ${this.tableName}`;
    }
}

export default base;
