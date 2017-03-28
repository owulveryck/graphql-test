import React, { Component } from 'react';
import './App.css';
import ReactTable from 'react-table'
import 'react-table/react-table.css'
import { ApolloClient, ApolloProvider, createNetworkInterface } from 'react-apollo';
import { gql, graphql } from 'react-apollo';

const columns = [{
  header: 'SKU',
  accessor: 'sku' // String-based value accessors!
}, {
  header: 'Location',
  accessor: 'location',
  sortable: true,
  width: 400,
  style: {
    textAlign: 'left'
  }
}, {
  header: 'Instance Type',
  accessor: 'instanceType'
}, {
  header: 'Operating System',
  accessor: 'operatingSystem'
}]

const offerColumns = [{
  header: 'Type',
  accessor: 'type',
  sortable: true,
},{
  header: 'Offering Class',
  accessor: 'OfferingClass',
} , {
  header: 'Purchase Option',
  accessor: 'PurchaseOption',
} , {
  header: 'Lease Contract Length',
  accessor: 'LeaseContractLength',

}]

const priceColumns = [{
  header: 'Unit',
  accessor: 'Unit',
},{
  header: 'Price',
  accessor: 'PricePerUnit',
}]


function ProductList({ loading, products }) {
  if (loading) {
    return <div>Loading</div>;
  } else {
    return (
      <div className="App">
      <ReactTable className="-striped -highlight"
      pivotBy={['location', 'instanceType']}
      data={products}
      columns={columns}
      SubComponent={(row) => {
        return (
          <div style={{padding: '20px'}}>
          <em>Offers</em>
          <br />
          <br />
          <ReactTable
          data={row.row.offers}
          columns={offerColumns}
          defaultPageSize={3}
          showPagination={true}
          SubComponent={(row) => {
            return (
              <div style={{padding: '20px'}}>
              <em>Prices</em>
              <br />
              <br />
              <ReactTable
              data={row.row.prices}
              columns={priceColumns}
              showPagination={true}
              defaultPageSize={3}
              />
              </div>
            );
          }}
            />
            </div>
          );
        }}
      />
      </div>
    );
  }
}

const allProducts = gql`
query products {
  products{
    sku
    location
    instanceType
    operatingSystem
    offers {
      type 
      OfferingClass
      LeaseContractLength
      PurchaseOption
      prices {
        Unit
        PricePerUnit
      }
    }
  }
}
`

const ProductListWithData = graphql(allProducts, {
  props: ({data: { loading, products }}) => ({
    loading,
      products,
  }),
})(ProductList);

class App extends Component {
  constructor(props) {
    super(props);
    const networkInterface = createNetworkInterface({
      uri: 'http://localhost:8080/graphql'
    })

    this.client = new ApolloClient({
      networkInterface: networkInterface
    });
  }

  render() {
    return (
      <ApolloProvider client={this.client}>
        <ProductListWithData />
      </ApolloProvider>
    );
  }
}

export default App;
