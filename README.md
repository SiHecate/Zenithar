# Zenithar

Proje adını, iş, ticaret, çalışma ve ticaretle ilişkilendirilen bir konsepti yansıtması için Skyrim oyunundan ilham alarak "Zenithar" olarak belirledim.

## Proje amacı

Net internet stajı için geliştirdiğim temel bir sipariş verme uygulaması.

## Açıklama
Backend development üzerine çalıştığım için projenin ağır kısmı backend üzerine, frontend konusunda temel işlevleri görecek şekilde bıraktım. CSS ve benzeri teknolojilerdeki zorlukları aşmak için ChatGPT gibi araçlardan yardım aldım, bu sayede frontend tecrübe eksikliğimi kapatmış oldum.

Oluşturduğum uygulama, ileriye dönük geliştirmelere açık bir yapıda tasarlandı ve kod tabanı, istenilen eklemeler veya çıkarmalar için kolaylıkla adapte edilebilir. Bu sayede projeyi sürekli olarak güncel tutmak, yeni özellikler eklemek veya gerektiğinde mevcut yapıyı değiştirmek oldukça kolaydır.

MVC (Model-View-Controller) yapısını temel aldım ve backend tarafında bu yapıyı en etkili şekilde kullanmaya odaklandım. Bu sayede kodun düzenli, sürdürülebilir ve genişletilebilir olmasını sağladım.

# Proje Açıklaması

Proje genel olarak kafe veya benzeri mekanlarda kullanılabilen bir sipariş yönetim uygulamasını içermektedir.

### Fonksiyonlar

- CRUD (Create - Read - Update - Delete) işlemleri yapıldı

- **Masa İşlemleri:**
  - Masa ekleme, silme, güncelleme ve listeleme işlemleri gerçekleştirilebilir.

- **Ürün İşlemleri:**
  - Ürün ekleme, silme, güncelleme ve listeleme işlemleri uygulanabilir.

- **Sipariş Alma ve Ödeme:**
  - Sipariş alma ve ödeme seçenekleri kullanıcıya sunulmuştur.

- **Güvenlik Önlemleri:**
  - Kullanıcı kimlik doğrulama sistemi eklenmiştir. Sadece admin yetkisine sahip kullanıcılar, belirli işlemleri gerçekleştirebilir.

- **Kullanıcı Yetkileri:**
  - Admin haricindeki kullanıcılar belirli işlevleri gerçekleştiremezler, bu da sistem

- **Monitoring ve Loglama Detayları:**
  - **Prometheus ve Grafana**: Proje performansını ve kaynak kullanımını izlemek için kullanılmıştır.
  - **Loki**: Logları toplamak, saklamak ve sorgulamak için kullanılmıştır.

Bu araçlar sayesinde proje üzerindeki performans, hata analizi ve sistem durumu hakkında detaylı bilgi elde edilebilir.


# Kullanılan Teknolojiler

Bu projede aşağıdaki teknolojiler kullanılmıştır:

- [Go](https://golang.org/) - Programlama Dili
- [Fiber](https://gofiber.io/) - Web Framework
- [GORM](https://gorm.io/) - ORM (Object-Relational Mapping) kütüphanesi
- [Docker](https://www.docker.com/) - Konteynerleşme Platformu
- [Fiber](https://gofiber.io/) - HTTP Framework
- [Prometheus](https://prometheus.io/) - İzleme ve Uyarı Sistemi
- [NGINX](https://www.nginx.com/) - Web Sunucusu
- [Loki](https://grafana.com/loki) - Log Aggregation Sistemi


