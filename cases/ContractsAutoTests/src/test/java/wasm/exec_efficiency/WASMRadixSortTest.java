package wasm.exec_efficiency;

import com.platon.rlp.datatypes.Int32;
import com.platon.rlp.datatypes.Int64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.RadixSort;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;
import java.util.Arrays;

/**
 * @title WASMRadixSortTest
 * @description 执行效率 - 基数排序
 * @author liweic
 * @updateTime 2020/3/3 17:09
 */
public class WASMRadixSortTest extends WASMContractPrepareTest {
    private String contractAddress;

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "wasm.exec_efficiency-基数排序", sourcePrefix = "wasm")
    public void test() {
        prepare();
        try {

            Integer numberOfCalls = Integer.valueOf(driverService.param.get("numberOfCalls"));

            RadixSort radixsort = RadixSort.load("0x1299ef235623d8471a3ce9a677515f9fa79b5386",web3j, transactionManager, provider);
            contractAddress = radixsort.getContractAddress();

            Int64[] arr = new Int64[numberOfCalls];

            int min = 1000, max = 4000;

            for (int i = 0; i < numberOfCalls; i++) {
                arr[i] = Int64.of(min + (int) (Math.random() * (max - min + 1)));
            }

            collector.logStepPass("before sort:" + Arrays.toString(arr));
            TransactionReceipt transactionReceipt = radixsort.load(contractAddress, web3j, transactionManager, provider)
                    .sort(arr, Int32.of(arr.length)).send();

            BigInteger gasUsed = transactionReceipt.getGasUsed();
            collector.logStepPass("gasUsed:" + gasUsed);
            collector.logStepPass("contract load successful. transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass("currentBlockNumber:" + transactionReceipt.getBlockNumber());

            Int64[] generationArr = radixsort.load(contractAddress, web3j, transactionManager, provider).get_array().send();

            collector.logStepPass("after sort:" + Arrays.toString(generationArr));
        } catch (Exception e) {
            e.printStackTrace();
            collector.logStepFail("The contract fail.", e.toString());
        }
    }

}