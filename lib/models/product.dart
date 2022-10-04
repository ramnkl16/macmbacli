import 'package:realm/realm.dart';
part 'product.g.dart';

@RealmModel()
class _Product {
  @PrimaryKey()
  @MapTo('_id')
  late ObjectId id;
  late String prodNum;
  late String prodType;
  late String servType;
  late String searchName; //strip out special chars and lower case
  late String name;
  late bool isKit;
  //late List<KV> attrs; //shows more info on product page
  // late List<_KV> specAttrs; //used for facet based filter
  // late List<_KV> typeAttrs; //Accessory Part,
}

class KV {
  late String k;
  late String v;
}
