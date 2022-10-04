// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'product.dart';

// **************************************************************************
// RealmObjectGenerator
// **************************************************************************

class Product extends _Product with RealmEntity, RealmObject {
  Product(
    ObjectId id,
    String prodNum,
    String prodType,
    String servType,
    String searchName,
    String name,
    bool isKit,
  ) {
    RealmObject.set(this, '_id', id);
    RealmObject.set(this, 'prodNum', prodNum);
    RealmObject.set(this, 'prodType', prodType);
    RealmObject.set(this, 'servType', servType);
    RealmObject.set(this, 'searchName', searchName);
    RealmObject.set(this, 'name', name);
    RealmObject.set(this, 'isKit', isKit);
  }

  Product._();

  @override
  ObjectId get id => RealmObject.get<ObjectId>(this, '_id') as ObjectId;
  @override
  set id(ObjectId value) => throw RealmUnsupportedSetError();

  @override
  String get prodNum => RealmObject.get<String>(this, 'prodNum') as String;
  @override
  set prodNum(String value) => RealmObject.set(this, 'prodNum', value);

  @override
  String get prodType => RealmObject.get<String>(this, 'prodType') as String;
  @override
  set prodType(String value) => RealmObject.set(this, 'prodType', value);

  @override
  String get servType => RealmObject.get<String>(this, 'servType') as String;
  @override
  set servType(String value) => RealmObject.set(this, 'servType', value);

  @override
  String get searchName =>
      RealmObject.get<String>(this, 'searchName') as String;
  @override
  set searchName(String value) => RealmObject.set(this, 'searchName', value);

  @override
  String get name => RealmObject.get<String>(this, 'name') as String;
  @override
  set name(String value) => RealmObject.set(this, 'name', value);

  @override
  bool get isKit => RealmObject.get<bool>(this, 'isKit') as bool;
  @override
  set isKit(bool value) => RealmObject.set(this, 'isKit', value);

  @override
  Stream<RealmObjectChanges<Product>> get changes =>
      RealmObject.getChanges<Product>(this);

  static SchemaObject get schema => _schema ??= _initSchema();
  static SchemaObject? _schema;
  static SchemaObject _initSchema() {
    RealmObject.registerFactory(Product._);
    return const SchemaObject(Product, 'Product', [
      SchemaProperty('_id', RealmPropertyType.objectid,
          mapTo: '_id', primaryKey: true),
      SchemaProperty('prodNum', RealmPropertyType.string),
      SchemaProperty('prodType', RealmPropertyType.string),
      SchemaProperty('servType', RealmPropertyType.string),
      SchemaProperty('searchName', RealmPropertyType.string),
      SchemaProperty('name', RealmPropertyType.string),
      SchemaProperty('isKit', RealmPropertyType.bool),
    ]);
  }
}
