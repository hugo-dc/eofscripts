yulcode=$1
#cat valid_ops_creator.yul
creator_bin=$(yul_comp "$(cat $yulcode)" | tail -n1)
#echo $creator_bin

creator_initcode_yul=$(./build/yulreturn $creator_bin)
#echo $creator_initcode_yul
creator_initcode_bin=$(yul_comp "$creator_initcode_yul" | tail -n1)
echo $creator_initcode_bin

./build/deploy "$creator_initcode_bin"
