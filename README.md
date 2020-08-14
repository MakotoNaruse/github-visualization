# gtihub-visualization

githubのリポジトリデータからエンジニアの生産性を可視化するツールを開発中

## 実装中めも

### 出す指標
- 言語毎の能力・経験
- それぞれの言語でどれほどのパフォーマンスを出しているか
- それを３段階の指標で測る
  - 個人でのスコア(コミットでの単純なコントリビュート)
  - 他者へのフィードバックのスコア(レビューでのコントリビュート)
  - 社会への貢献のスコア(OSSへのコントリビュート)

### 課題
- githubのリポジトリ・プロフィールページを一目見ただけではその人がどの分野でどれほどの能力・経験があるのかが分からない
- 一目見て能力・経験の広さ/深さを見たい
- それに加えて、時系列的、加速度的な経験(能力)の増加も見たい
- 段階的な指標を用いることで、自分の到達レベルを知りたい
- 明確に到達レベルを分けることで他社からのフィードバックをわかりやすくする
  - 経験が足りないよっていうフィードバックに悩まされた自身の経験から
  - もっと他者への貢献（チーム開発の経験）も欲しいよねといったフォードバック出せるようにする
- 自分自身が、他者へのフィードバックや社会への貢献がエンジニアの高い能力を表すと考えている
  - 自分の考え
    - そもそもエンジニアリングは１人でやるものではない
    - 今自分が学んでいることも誰かのアウトプットの結果
    - 成長したエンジニアは他者・社会へのアウトプットをしてエンジニアリング全体の能力向上を計る

### 課題解決・価値としての落とし込み方
現段階
- 段階別・時系列でのグラフ表示（能力の可視化）
---
展望
- 未来への提案（ギャップの可視化）
  - ロールモデルの提示
- 指標を積むためのマッチングの提案、機会の提供（ギャップを埋める手法）
